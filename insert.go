package chorm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/0x19/go-clickhouse-model/dml"
	"github.com/0x19/go-clickhouse-model/models"
	"github.com/vahid-sohrabloo/chconn/v2"
	"github.com/vahid-sohrabloo/chconn/v2/column"
)

type InsertBuilder[T models.Model] struct {
	ctx     context.Context
	orm     *ORM
	model   T
	builder *dml.InsertBuilder
}

func (b *InsertBuilder[T]) Build() (string, error) {
	return b.builder.Build()
}

func (b *InsertBuilder[T]) ExecContext(ctx context.Context, queryOptions *chconn.QueryOptions, columns ...column.ColumnBasic) error {
	return b.orm.GetConn().InsertWithOption(ctx, b.SQL(), queryOptions, columns...)
}

func (b *InsertBuilder[T]) SQL() string {
	return b.builder.String()
}

func NewInsert[T models.Model](ctx context.Context, orm *ORM, model T, queryOptions *chconn.QueryOptions) (T, *InsertBuilder[T], error) {
	// Check if the underlying value of the interface is nil. Unfortunately, it is a T and we cannot
	// directly check if it's nil due to type missmatch.
	{
		modelValue := reflect.ValueOf(model)

		if !modelValue.IsValid() {
			return model, nil, fmt.Errorf("underlying value of model cannot be nil")
		}

		if modelValue.Kind() == reflect.Ptr && modelValue.IsNil() {
			return model, nil, fmt.Errorf("underlying value of model cannot be nil")
		}
	}

	stmtBuilder := dml.NewInsertBuilder()
	stmtBuilder.Model(model)
	stmtBuilder.Fields(GetModelKeys(model)...)

	builder := &InsertBuilder[T]{
		ctx:     ctx,
		orm:     orm,
		model:   model,
		builder: stmtBuilder,
	}

	if err := builder.ExecContext(ctx, queryOptions); err != nil {
		return model, builder, err
	}

	return model, builder, nil
}
