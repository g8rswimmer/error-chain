package chain

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestErrorChain_Errors(t *testing.T) {
	type fields struct {
		head *link
		tail *link
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "errors",
			fields: fields{
				head: &link{
					err: errors.New("1"),
					next: &link{
						err: errors.New("2"),
						next: &link{
							err: errors.New("3"),
						},
					},
				},
			},
			want: []error{
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorChain{
				head: tt.fields.head,
				tail: tt.fields.tail,
			}
			if got := e.Errors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorChain.Errors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorChain_Add(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		ec   *ErrorChain
		args args
		want []error
	}{
		{
			name: "From Nothing",
			ec: func() *ErrorChain {
				return &ErrorChain{}
			}(),
			args: args{
				err: errors.New("1"),
			},
			want: []error{
				errors.New("1"),
			},
		},
		{
			name: "From Something",
			ec: func() *ErrorChain {
				e := &ErrorChain{
					head: &link{
						err: errors.New("1"),
					},
				}
				e.tail = e.head
				return e
			}(),
			args: args{
				err: errors.New("2"),
			},
			want: []error{
				errors.New("1"),
				errors.New("2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ec.Add(tt.args.err)
			if got := tt.ec.Errors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorChain.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorChain_Error(t *testing.T) {
	type fields struct {
		head *link
		tail *link
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string",
			fields: fields{
				head: &link{
					err: errors.New("1"),
					next: &link{
						err: errors.New("2"),
						next: &link{
							err: errors.New("3"),
						},
					},
				},
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorChain{
				head: tt.fields.head,
				tail: tt.fields.tail,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ErrorChain.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testError struct {
	code int
}

func (t *testError) Error() string {
	return fmt.Sprintf("%d", t.code)
}

func (t *testError) Is(target error) bool {
	te, ok := target.(*testError)
	if ok == false {
		return false
	}
	return t.code == te.code
}

func TestErrorChain_Unwrap(t *testing.T) {

	type fields struct {
		head *link
		tail *link
	}
	tests := []struct {
		name   string
		fields fields
		target error
	}{
		{
			name: "unwrap it",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
					next: &link{
						err: &testError{
							code: 2,
						},
						next: &link{
							err: &testError{
								code: 3,
							},
						},
					},
				},
			},
			target: &testError{
				code: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorChain{
				head: tt.fields.head,
				tail: tt.fields.tail,
			}
			err := e.Unwrap()
			if errors.Is(err, tt.target) == false {
				t.Errorf("ErrorChain.Unwrap() = %v, want %v", err, tt.target)
			}
		})
	}
}

func TestErrorChain_Is(t *testing.T) {
	type fields struct {
		head *link
		tail *link
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Simple",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
				},
			},
			args: args{
				target: &testError{
					code: 1,
				},
			},
			want: true,
		},
		{
			name: "Simple no found",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
				},
			},
			args: args{
				target: &testError{
					code: 2,
				},
			},
			want: false,
		},
		{
			name: "Simple wrapped",
			fields: fields{
				head: &link{
					err: func() error {
						te := &testError{
							code: 1,
						}
						return fmt.Errorf("test wrap %w", te)
					}(),
				},
			},
			args: args{
				target: &testError{
					code: 1,
				},
			},
			want: true,
		},
		{
			name: "Simple wrapped not found",
			fields: fields{
				head: &link{
					err: func() error {
						te := &testError{
							code: 1,
						}
						return fmt.Errorf("test wrap %w", te)
					}(),
				},
			},
			args: args{
				target: &testError{
					code: 2,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorChain{
				head: tt.fields.head,
				tail: tt.fields.tail,
			}
			if got := e.Is(tt.args.target); got != tt.want {
				t.Errorf("ErrorChain.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorChain_ErrorsIs(t *testing.T) {
	type fields struct {
		head *link
		tail *link
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Simple",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
				},
			},
			args: args{
				target: &testError{
					code: 1,
				},
			},
			want: true,
		},
		{
			name: "Simple no found",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
				},
			},
			args: args{
				target: &testError{
					code: 2,
				},
			},
			want: false,
		},
		{
			name: "Simple wrapped",
			fields: fields{
				head: &link{
					err: func() error {
						te := &testError{
							code: 1,
						}
						return fmt.Errorf("test wrap %w", te)
					}(),
				},
			},
			args: args{
				target: &testError{
					code: 1,
				},
			},
			want: true,
		},
		{
			name: "Simple wrapped not found",
			fields: fields{
				head: &link{
					err: func() error {
						te := &testError{
							code: 1,
						}
						return fmt.Errorf("test wrap %w", te)
					}(),
				},
			},
			args: args{
				target: &testError{
					code: 2,
				},
			},
			want: false,
		},
		{
			name: "Complex",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
					next: &link{
						err: &testError{
							code: 2,
						},
					},
				},
			},
			args: args{
				target: &testError{
					code: 2,
				},
			},
			want: true,
		},
		{
			name: "Complex Not found",
			fields: fields{
				head: &link{
					err: &testError{
						code: 1,
					},
					next: &link{
						err: &testError{
							code: 2,
						},
					},
				},
			},
			args: args{
				target: &testError{
					code: 22,
				},
			},
			want: false,
		},
		{
			name: "Complex wrapped",
			fields: fields{
				head: &link{
					err: func() error {
						te := &testError{
							code: 1,
						}
						return fmt.Errorf("test wrap %w", te)
					}(),
					next: &link{
						err: func() error {
							te := &testError{
								code: 22,
							}
							return fmt.Errorf("test wrap %w", te)
						}(),
					},
				},
			},
			args: args{
				target: &testError{
					code: 22,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorChain{
				head: tt.fields.head,
				tail: tt.fields.tail,
			}
			if got := errors.Is(e, tt.args.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
