module Sessions::Operation
  class Start < Trailblazer::Operation
    step :find_user
    step :authenticate
    step :create_session

    def find_user(ctx, params:, **)
      ctx["user"] = User.find_by(id: params[:user_id])
    end

    def authenticate(ctx, params:, **)
      ctx["user"].authenticate(params[:password])
    end

    def create_session(ctx, **)
      ctx["session"] = ctx["user"].user_sessions.create!(
        user_agent: ctx["request"].user_agent,
        ip_address: ctx["request"].remote_addr
      )
    end
  end
end
