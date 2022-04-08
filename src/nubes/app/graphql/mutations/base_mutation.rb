module Mutations
  class BaseMutation < GraphQL::Schema::RelayClassicMutation
    argument_class Types::Base::Argument
    field_class Types::Base::Field
    input_object_class Types::Base::InputObject
    object_class Types::Base::Object

    protected

    def to_errors(error)
      return [{ message: error, code: error }] unless error.is_a?(ActiveModel::Errors)

      error.map do |err|
        {
          path: ["input", err.attribute.to_s.camelize(:lower)],
          message: err.message,
          code: err.type
        }
      end
    end
  end
end
