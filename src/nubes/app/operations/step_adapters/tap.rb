module StepAdapters
  # Add a "tap" step that works similarly to the built-in "tee" step in
  # that if the step is successful it ignores the output and returns the
  # input.  However, if the step returns a Failure, then that failure is
  # not ignored as it is with "tee", but is instead returned directly.
  class Tap
    include Dry::Monads::Result::Mixin

    def call(operation, _options, args)
      result = operation.call(*args)

      if result.respond_to?(:failure?) && result.failure?
        return result
      end

      Success(args[0])
    end
  end
end
