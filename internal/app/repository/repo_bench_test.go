package repository

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"

	"GRPCService/config"
	"GRPCService/internal/models"

	"github.com/google/uuid"
)

func getInternalStorage() Repository {
	return NewInnerStorageRepository()
}

func BenchmarkInternalMemory(b *testing.B) {
	ctx := context.Background()

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Println(err)
		}
	}()

	repo := getInternalStorage()

	b.ResetTimer()

	b.Run("BenchmarkInternalInsert", benchmarkInsert(ctx, repo))
	b.Run("BenchmarkInternalGet", benchmarkGet(ctx, repo))
	b.Run("BenchmarkInternalUpd", benchmarkUpd(ctx, repo))
	b.Run("BenchmarkInternalDelete", benchmarkDelete(ctx, repo))
}

func getScylla(ctx context.Context) (Repository, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	repo, err := NewScyllaRepository(ctx, conf.ScyllaAddr)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func BenchmarkScylla(b *testing.B) {
	ctx := context.Background()

	repo, err := getScylla(ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("BenchmarkScyllaInsert", benchmarkInsert(ctx, repo))
	b.Run("BenchmarkScyllaGet", benchmarkGet(ctx, repo))
	b.Run("BenchmarkScyllaUpd", benchmarkUpd(ctx, repo))
	b.Run("BenchmarkScyllaDelete", benchmarkDelete(ctx, repo))
}

func getScyllaX(ctx context.Context) (Repository, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	repo, err := NewScyllaXRepository(ctx, conf.ScyllaAddr)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func BenchmarkCQLX(b *testing.B) {
	ctx := context.Background()

	repo, err := getScyllaX(ctx)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("BenchmarkScyllaXInsert", benchmarkInsert(ctx, repo))
	b.Run("BenchmarkScyllaXGet", benchmarkGet(ctx, repo))
	b.Run("BenchmarkScyllaXUpd", benchmarkUpd(ctx, repo))
	b.Run("BenchmarkScyllaXDelete", benchmarkDelete(ctx, repo))
}

func getMemcached() (Repository, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	repo, err := NewMemcacheRepository(conf.MCServerAddr)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func BenchmarkMemcached(b *testing.B) {
	ctx := context.Background()

	repo, err := getMemcached()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("BenchmarkMemcachedInsert", benchmarkInsert(ctx, repo))
	b.Run("BenchmarkMemcachedUpd", benchmarkUpd(ctx, repo))
	b.Run("BenchmarkMemcachedGet", benchmarkGet(ctx, repo))
	b.Run("BenchmarkMemcachedDelete", benchmarkDelete(ctx, repo))
}

func benchmarkInsert(ctx context.Context, repo Repository) func(b *testing.B) {
	r1 := models.Record{
		ID:   uuid.Nil,
		Data: "message 1",
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := repo.SetRecord(ctx, r1)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchmarkUpd(ctx context.Context, repo Repository) func(b *testing.B) {
	r1 := models.Record{
		ID:   uuid.Nil,
		Data: "message 1",
	}
	r2 := models.Record{
		Data: "message 2",
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()

			r, err := repo.SetRecord(ctx, r1)
			if err != nil {
				b.Fatal(err)
			}

			r2.ID = r.ID

			b.StartTimer()

			_, err = repo.SetRecord(ctx, r2)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchmarkGet(ctx context.Context, repo Repository) func(b *testing.B) {
	r1 := models.Record{
		ID:   uuid.Nil,
		Data: "message 1",
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()

			r, err := repo.SetRecord(ctx, r1)
			if err != nil {
				b.Fatal(err)
			}

			b.StartTimer()

			_, err = repo.GetRecord(ctx, r.ID)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchmarkDelete(ctx context.Context, repo Repository) func(b *testing.B) {
	r1 := models.Record{
		ID:   uuid.Nil,
		Data: "message 1",
	}

	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()

			r, err := repo.SetRecord(ctx, r1)
			if err != nil {
				b.Fatal(err)
			}

			b.StartTimer()

			_, err = repo.DeleteRecord(ctx, r.ID)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
