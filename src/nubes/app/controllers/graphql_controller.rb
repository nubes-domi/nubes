class GraphqlController < ApplicationController
  skip_before_action :ensure_authenticated

  # If accessing from outside this domain, nullify the session
  # This allows for outside API access while preventing CSRF attacks,
  # but you'll have to authenticate your user separately
  protect_from_forgery with: :null_session

  def execute
    variables = prepare_variables(params[:variables])
    query = params[:query]
    operation_name = params[:operationName]
    result = NubesSchema.execute(query, variables:, context: { current_user: }, operation_name:)
    render json: result
  end

  private

  # Handle variables in form data, JSON body, or a blank value
  def prepare_variables(variables_param)
    return {} unless variables_param.present?
    return variables_param if variables_param.is_a?(Hash)
    return JSON.parse(variables_param) || {} if variables_param.is_a?(String)

    # GraphQL-Ruby will validate name and type of incoming variables.
    return variables_param.to_unsafe_hash if variables_param.is_a?(ActionController::Parameters)

    raise ArgumentError, "Unexpected parameter: #{variables_param}"
  end
end
