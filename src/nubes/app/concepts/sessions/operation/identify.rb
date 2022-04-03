module Sessions
  module Operation
    class Identify < Trailblazer::Operation
      step :identify

      def identify(ctx, params:, **)
        ctx["user"] = User.identify(params[:identifier])
      end
    end
  end
end
