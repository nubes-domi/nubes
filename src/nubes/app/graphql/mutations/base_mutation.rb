module Mutations
  class BaseMutation < GraphQL::Schema::RelayClassicMutation
    include Dry::Monads[:result]

    argument_class Types::Base::Argument
    field_class Types::Base::Field
    input_object_class Types::Base::InputObject
    object_class Types::Base::Object

    protected

    def handle_failures(result, &success)
      case result
      in Success(result)
        success.call(result)
      in Failure(Dry::Schema::Result => schema)
        { errors: to_schema_errors(schema.errors) }
      in Failure(ArInvalid => errors)
        { errors: to_errors(errors) }
      in Failure(NotFound | Forbidden)
        { errors: [{ message: "not_found" }] }
      end
    end

    def to_schema_errors(errors)
      errors.map do |err|
        path = err.path
        path = path[1..] if path[0] == :attributes

        {
          path: ["input", path[0].to_s.camelize(:lower)],
          message: err.text,
          code: nil
        }
      end
    end

    def to_errors(errors)
      errors.map do |err|
        {
          path: ["input", err.attribute.to_s.camelize(:lower)],
          message: err.message,
          code: err.type
        }
      end
    end
  end
end
