class ApplicationController < ActionController::Base
  protected

  def start_session(response)
    cookies.permanent[:current_session] = response.access_token
    # cookies.permanent[:sessions] = [(cookies[:sessions] || ""), session.access_token].join("|")
  end

  before_action do
    Thread.current[:authorization] = cookies[:current_session]
  end

  after_action do
    Thread.current[:authorization] = nil
  end
end
