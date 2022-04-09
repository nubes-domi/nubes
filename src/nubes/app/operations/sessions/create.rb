module Sessions
  class Create < BaseOperation
    input_validation do
      required(:user_id).filled(:string)
      required(:password).filled(:string)
      required(:request).filled
    end

    fetch User, :user_id, merge_as: :user
    tap :authenticate
    merge :build
    step :create

    def authenticate(user:, password:, **)
      if user.authenticate(password)
        Success(user)
      else
        Failure(NotFound.new)
      end
    end

    def build(user:, request:, **)
      record = user.user_sessions.create!(
        user_agent: request.user_agent,
        ip_address: request.remote_addr
      )
      Success({ record: })
    end
  end
end
