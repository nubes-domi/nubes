module SessionsManagement
  extend ActiveSupport::Concern

  included do
    before_action :ensure_authenticated
  end

  protected

  def start_session(session)
    cookies[:current_session] = session.token if add_session(session)
  end

  def terminate_session(session)
    delete_session(session)
    cookies.delete(:current_session) if cookies[:current_session] == session.token
  end

  def active_sessions
    @active_sessions ||= begin
      sessions = (cookies[:sessions] || "").split("|")
      sessions.to_h { |token| [token, UserSession.for_token(token)] }.compact
    end
  end

  def switch_session(token)
    if active_sessions.key?(token)
      cookies[:current_session] = token
    else
      false
    end
  end

  def current_session
    return unless cookies[:current_session].present?

    session = UserSession.for_token(cookies[:current_session])
    cookies.delete(:current_session) unless session

    session
  end

  def current_user
    @current_user ||= current_session.user
  end

  def session_for?(user_id)
    !!session_for(user_id)
  end

  def session_for(user_id)
    active_sessions.values.find { |session| session.user_id == user_id }
  end

  def session_token_for(user_id)
    pair = active_sessions.find { |_, session| session.user_id == user_id }
    pair[0] if pair
  end

  def ensure_authenticated
    redirect_to signin_path(continue: request.path) unless current_session
  end

  private

  def add_session(session)
    sessions = active_sessions
    return false if session_for?(session.user_id)

    sessions[session.token] = session
    save_sessions_cookie(sessions)
    true
  end

  def delete_session(session)
    sessions = active_sessions
    sessions.delete(session.token)
    save_sessions_cookie(sessions)
  end

  def save_sessions_cookie(sessions)
    sessions.uniq { |_, session| session.user_id }
    cookies[:sessions] = sessions.keys.join("|")
  end
end
