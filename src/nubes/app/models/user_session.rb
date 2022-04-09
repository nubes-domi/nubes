class UserSession < ApplicationRecord
  include PrettyId
  self.id_prefix = :usr_sess

  belongs_to :user

  before_create do
    self.expires_at ||= 10.years.from_now
  end

  class << self
    def for_token(token)
      decoded = JWT.decode(token, Rails.application.secret_key_base, true, {
        algorithm: "HS256", aud: "auth", verify_aud: true
      })

      find_by(id: decoded[0]["jti"])
    rescue JWT::DecodeError, JWT::InvalidAudError, JWT::ExpiredSignature => e
      nil
    end
  end

  def token
    @token ||= JWT.encode({
      sub: user_id,
      jti: id,
      exp: expires_at.to_i,
      aud: "auth"
    }, Rails.application.secret_key_base, "HS256")
  end
end
