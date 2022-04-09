module StepAdapters
  # Adds a "merge" step that expects the input to the step to be a hash,
  # and then merges the successful result of the step into the input hash.
  # If the step results in a Failure then that Failure is returned instead.
  #
  # If the return value of the step is not a hash, then the return value is
  # merged into the input hash using the step name as the key.
  class Merge < Dry::Transaction::StepAdapters
    include Dry::Monads[:result]

    def call(operation, _options, args)
      result = operation.call(args[0])
      return result if result.failure?

      value = result.value!
      unless value.is_a?(Hash)
        value = { operation.operation.name => value }
      end

      Success(args[0].merge(value))
    end
  end
end
