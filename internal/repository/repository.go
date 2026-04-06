package repository

type Executor interface{}

type Repository struct {
	executor Executor
}

func NewRepository(executor Executor) *Repository {
	return &Repository{
		executor: executor,
	}
}
