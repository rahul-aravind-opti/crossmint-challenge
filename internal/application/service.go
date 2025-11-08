package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/crossmint/megaverse-challenge/internal/application/strategies"
	"github.com/crossmint/megaverse-challenge/internal/domain"
	"github.com/crossmint/megaverse-challenge/internal/domain/entities"
)

// MegaverseService orchestrates the creation and management of megaverses
type MegaverseService struct {
	repository domain.MegaverseRepository
	logger     *log.Logger
	limiter    rateLimiter
}

type rateLimiter interface {
	Wait(context.Context) error
}

// NewMegaverseService creates a new megaverse service
func NewMegaverseService(repository domain.MegaverseRepository, logger *log.Logger, limiter rateLimiter) *MegaverseService {
	if logger == nil {
		logger = log.Default()
	}
	return &MegaverseService{
		repository: repository,
		logger:     logger,
		limiter:    limiter,
	}
}

// ExecuteStrategy executes a pattern strategy to create a megaverse
func (s *MegaverseService) ExecuteStrategy(ctx context.Context, strategy strategies.PatternStrategy) error {
	s.logger.Printf("Executing strategy: %s\n", strategy.GetName())

	// Generate the plan
	plan, err := strategy.GeneratePlan(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate plan: %w", err)
	}

	objects := plan.Objects

	s.logger.Printf("Generated plan with %d objects\n", len(objects))

	// Create objects based on the execution order
	execOrder := plan.Order
	if execOrder == 0 {
		execOrder = strategies.OrderSequential
	}

	batchSize := plan.BatchSize
	if batchSize <= 0 {
		batchSize = 5
	}

	switch execOrder {
	case strategies.OrderParallel:
		return s.createObjectsParallel(ctx, objects)
	case strategies.OrderBatched:
		return s.createObjectsBatched(ctx, objects, batchSize)
	default:
		return s.createObjectsSequential(ctx, objects)
	}
}

// createObjectsSequential creates objects one by one
func (s *MegaverseService) createObjectsSequential(ctx context.Context, objects []entities.AstralObject) error {
	totalObjects := len(objects)
	var errs []error

	for i, obj := range objects {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		s.logger.Printf("[%d/%d] Creating %s at position (%d, %d)\n",
			i+1, totalObjects, obj.GetType(), obj.GetPosition().Row, obj.GetPosition().Column)

		if err := s.waitForRateLimit(ctx); err != nil {
			return err
		}

		if err := s.createObject(ctx, obj); err != nil {
			s.logger.Printf("Failed to create object at (%d, %d): %v\n",
				obj.GetPosition().Row, obj.GetPosition().Column, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during sequential creation", len(errs))
	}

	return nil
}

// createObjectsParallel creates objects in parallel with controlled concurrency
func (s *MegaverseService) createObjectsParallel(ctx context.Context, objects []entities.AstralObject) error {
	// Use a worker pool pattern with limited concurrency
	const maxWorkers = 5 // Limit concurrent workers to avoid overwhelming the API

	var wg sync.WaitGroup
	objectChan := make(chan entities.AstralObject, len(objects))
	errorChan := make(chan error, len(objects))

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for obj := range objectChan {
				if err := ctx.Err(); err != nil {
					errorChan <- err
					return
				}

				if err := s.waitForRateLimit(ctx); err != nil {
					s.logger.Printf("[Worker %d] rate limit wait failed: %v\n", workerID, err)
					errorChan <- err
					return
				}

				s.logger.Printf("[Worker %d] Creating %s at position (%d, %d)\n",
					workerID, obj.GetType(), obj.GetPosition().Row, obj.GetPosition().Column)

				if err := s.createObject(ctx, obj); err != nil {
					s.logger.Printf("[Worker %d] Failed to create object: %v\n", workerID, err)
					errorChan <- err
					// Continue processing other objects
				}
			}
		}(i)
	}

	// Send objects to workers
	for _, obj := range objects {
		objectChan <- obj
	}
	close(objectChan)

	// Wait for all workers to complete
	wg.Wait()
	close(errorChan)

	var errs []error
	for err := range errorChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during parallel creation", len(errs))
	}

	return nil
}

// createObjectsBatched creates objects in batches
func (s *MegaverseService) createObjectsBatched(ctx context.Context, objects []entities.AstralObject, batchSize int) error {
	totalObjects := len(objects)

	var errs []error

	for i := 0; i < totalObjects; i += batchSize {
		end := i + batchSize
		if end > totalObjects {
			end = totalObjects
		}

		batch := objects[i:end]
		s.logger.Printf("Processing batch %d-%d of %d\n", i+1, end, totalObjects)

		// Create batch sequentially (could be parallel within batch)
		for _, obj := range batch {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("context cancelled: %w", err)
			}

			if err := s.waitForRateLimit(ctx); err != nil {
				return err
			}

			if err := s.createObject(ctx, obj); err != nil {
				s.logger.Printf("Failed to create object: %v\n", err)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during batched creation", len(errs))
	}

	return nil
}

func (s *MegaverseService) waitForRateLimit(ctx context.Context) error {
	if s.limiter == nil {
		return nil
	}
	if err := s.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait failed: %w", err)
	}
	return nil
}

// createObject creates a single astral object using the repository
func (s *MegaverseService) createObject(ctx context.Context, obj entities.AstralObject) error {
	if err := obj.Validate(); err != nil {
		return fmt.Errorf("invalid object: %w", err)
	}

	pos := obj.GetPosition()

	switch o := obj.(type) {
	case *entities.Polyanet:
		return s.repository.CreatePolyanet(ctx, pos)

	case *entities.Soloon:
		return s.repository.CreateSoloon(ctx, pos, o.Color)

	case *entities.Cometh:
		return s.repository.CreateCometh(ctx, pos, o.Direction)

	default:
		return fmt.Errorf("unknown object type: %T", obj)
	}
}

// ClearMegaverse removes all objects from the megaverse
func (s *MegaverseService) ClearMegaverse(ctx context.Context, width, height int) error {
	s.logger.Printf("Clearing megaverse (%dx%d)\n", width, height)

	successCount := 0
	errorCount := 0

	// Try to delete objects at all positions
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			pos := entities.Position{Row: row, Column: col}
			deleted := false

			// Try deleting each type of object (we don't know what's there)
			// The API will return an error if nothing exists, which we can ignore
			for _, objType := range []string{"POLYANET", "SOLOON", "COMETH"} {
				err := s.repository.DeleteObject(ctx, objType, pos)
				if err == nil {
					successCount++
					s.logger.Printf("Deleted %s at (%d, %d)\n", objType, row, col)
					deleted = true
					break // Successfully deleted something, move to next position
				}
				// If error, try the next type
			}

			if !deleted {
				errorCount++
			}

			// Small delay to respect rate limits
			time.Sleep(100 * time.Millisecond)
		}
	}

	s.logger.Printf("Clear complete: %d objects removed, %d positions checked, %d positions unchanged\n",
		successCount, width*height, errorCount)

	return nil
}

// ValidateMegaverse validates the current megaverse against a strategy
func (s *MegaverseService) ValidateMegaverse(ctx context.Context, strategy strategies.PatternStrategy) error {
	s.logger.Printf("Validating megaverse against strategy: %s\n", strategy.GetName())

	// Try to get the current map
	currentMap, err := s.repository.GetCurrentMap(ctx)
	if err != nil {
		// If current map is not available, we can't validate
		s.logger.Printf("Warning: Could not retrieve current map: %v\n", err)
		s.logger.Printf("Validation skipped - current map endpoint may not be available\n")
		return nil
	}

	// Validate against the strategy
	if err := strategy.Validate(currentMap); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Printf("Validation successful!\n")
	return nil
}

// GetGoalMap retrieves the goal map for the current challenge
func (s *MegaverseService) GetGoalMap(ctx context.Context) (*domain.GoalMap, error) {
	return s.repository.GetGoalMap(ctx)
}
