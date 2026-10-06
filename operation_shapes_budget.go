package apitools

import (
	"context"
	"encoding/json"

	"github.com/OpenUdon/uws/binding"
)

// One invocation owns its counters. Bounds apply before accumulating a whole
// source or compiling another reused schema; hard process isolation is external.
type operationShapeBudget struct {
	ctx           context.Context
	maxOperations int
	operations    int
	tableBytes    int
	schemaBytes   int
	schemaWork    int
	schemaDialect int
	err           error
}

const (
	shapeSchemaDraft7 = iota
	shapeSchemaOpenAPI30
	shapeSchema2020
)

func newOperationShapeBudget(ctx context.Context, maxOperations int) *operationShapeBudget {
	// Include the empty envelope and one conservative separator per record.
	data, _ := json.Marshal(binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{}, Operations: []binding.OperationShape{}})
	return &operationShapeBudget{ctx: ctx, maxOperations: maxOperations, tableBytes: len(data)}
}

func (budget *operationShapeBudget) source(source binding.Source) error {
	if err := budget.ctx.Err(); err != nil {
		return err
	}
	if (binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{source}}).Validate() != nil {
		return ErrOperationShapeTable
	}
	data, err := json.Marshal(source)
	if err != nil || len(data)+1 > binding.MaxTableBytes-budget.tableBytes {
		return ErrOperationShapeTable
	}
	budget.tableBytes += len(data) + 1
	return nil
}

func (budget *operationShapeBudget) append(operations *[]binding.OperationShape, shape binding.OperationShape) error {
	if err := budget.ctx.Err(); err != nil {
		return err
	}
	if budget.err != nil {
		return budget.err
	}
	if budget.operations >= budget.maxOperations || len(shape.Inputs) > 2048 || len(shape.Outputs) > 2048 {
		return ErrOperationShapeTable
	}
	// Refuse amplified schemas before the UWS structural validator/marshaler
	// can traverse or allocate their full expanded representation.
	remaining := binding.MaxTableBytes - budget.tableBytes
	for _, input := range shape.Inputs {
		if len(input.Schema.JSON) > remaining {
			return ErrOperationShapeTable
		}
		remaining -= len(input.Schema.JSON)
	}
	for _, output := range shape.Outputs {
		if len(output.Schema.JSON) > remaining {
			return ErrOperationShapeTable
		}
		remaining -= len(output.Schema.JSON)
	}
	if (binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{shape.Source}, Operations: []binding.OperationShape{shape}}).Validate() != nil {
		return ErrOperationShapeTable
	}
	data, err := json.Marshal(shape)
	if err != nil || len(data)+1 > binding.MaxTableBytes-budget.tableBytes {
		return ErrOperationShapeTable
	}
	budget.tableBytes += len(data) + 1
	budget.operations++
	if operations != nil {
		*operations = append(*operations, shape)
	}
	return nil
}

func boundedShapeTableBytes(ctx context.Context, table binding.ShapeTable) ([]byte, error) {
	if len(table.Sources) > binding.MaxSources || len(table.Operations) > binding.MaxOperations {
		return nil, ErrOperationShapeTable
	}
	budget := newOperationShapeBudget(ctx, binding.MaxOperations)
	for _, source := range table.Sources {
		if err := budget.source(source); err != nil {
			return nil, err
		}
	}
	for _, shape := range table.Operations {
		if err := budget.append(nil, shape); err != nil {
			return nil, err
		}
	}
	return table.Marshal()
}
