Dry::Transaction::StepAdapters.register(:merge, StepAdapters::Merge.new)
Dry::Transaction::StepAdapters.register(:tap, StepAdapters::Tap.new)
Dry::Transaction::StepAdapters.register(:valid, StepAdapters::Validate.new)

Forbidden = Class.new(StandardError)
Invalid = Class.new(StandardError)
ArInvalid = Class.new(StandardError)
NotFound = Class.new(StandardError)

class BaseOperation
  include Dry::Transaction
  include Dry::Monads[:maybe, :result, :do]
  include ::StepAdapters::Validate::Mixin

  class << self
    def call(**kwargs)
      new.call(**kwargs)
    end

    def fetch(type, id_key = :id, merge_as: :record)
      merge :fetch

      define_method :fetch do |params|
        id = params[id_key]
        object = type.find_by(id:)
        if object
          Success({ merge_as => object })
        else
          Failure(NotFound.new(type:, id:))
        end
      end
    end

    def authorize(key = :record, query = nil)
      step :authorize

      define_method :authorize do |params|
        record = key.is_a?(Class) ? key : params[key]
        query ||= "#{self.class.name.demodulize.underscore}?"
        policy = Pundit.policy!(params[:user], record)

        if policy.public_send(query)
          Success(params)
        else
          Failure(Forbidden.new(query:, record:, policy:))
        end
      end
    end
  end

  def create(record:, **)
    if record.save
      Success(record)
    else
      Failure(ArInvalid.new(record.errors))
    end
  end

  def update(record:, attributes:, **)
    if record.update(attributes)
      Success(record)
    else
      Failure(ArInvalid.new(record.errors))
    end
  end

  def destroy(record:, **)
    if record.destroy
      Success(true)
    else
      Failure(ArInvalid.new(record.errors))
    end
  end

  def transaction(**any, &block)
    result = nil
    ActiveRecord::Base.transaction do
      result = block.call(Success(**any))
      raise ActiveRecord::Rollback if result.failure?
    end
    result
  rescue ActiveRecord::Rollback
    Failure(result)
  end
end
