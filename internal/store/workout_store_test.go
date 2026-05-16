package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*sql.DB, int) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run the migrations for our test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE users, workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables %v", err)
	}

	var userID int
	err = db.QueryRow(`INSERT INTO users (username, email, password_hash, bio) VALUES ($1, $2, $3, $4) RETURNING id`, "test_user", "test@example.com", "hash", "bio").Scan(&userID)
	if err != nil {
		t.Fatalf("inserting test user: %v", err)
	}

	return db, userID
}

func TestCreateWorkout(t *testing.T) {
	db, userID := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				UserID:          userID,
				Title:           "pushup day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "bench press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(153.5),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid workout",
			workout: &Workout{
				UserID:          99999, // Non-existent UserID to trigger foreign key violation error
				Title:           "full body",
				Description:     "complete body workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "squat",
						Sets:         4,
						Reps:         IntPtr(60),
						Weight:       FloatPtr(153.5),
						Notes:        "full depth",
						OrderIndex:   1,
					}, {
						ExerciseName:    "plank",
						Sets:            4,
						Reps:            IntPtr(60),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(153.5),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
				ID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrived, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrived.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrived.Entries))

			for i := range retrived.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrived.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrived.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrived.Entries[i].OrderIndex)
			}
		})
	}
}


func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
