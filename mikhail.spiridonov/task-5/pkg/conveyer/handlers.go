package conveyer

import "context"

type DecoratorFunc func(

) error

type MultiplexerFunc func(

) error

type SeparatorFunc func(

) error