class ApplicationController < ActionController::Base
  protected

  def start_session(response)
    cookies.permanent[:current_session] = response.access_token
    # cookies.permanent[:sessions] = [(cookies[:sessions] || ""), session.access_token].join("|")
  end
end
