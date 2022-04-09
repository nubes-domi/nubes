module Sessions
  class Identify < BaseOperation
    input_validation do
      required(:identifier).filled(:string)
    end

    step :identify

    def identify(identifier:, **)
      user = User.identify(identifier)
      if user
        Success(user)
      else
        Failure(NotFound.new)
      end
    end
  end
end
