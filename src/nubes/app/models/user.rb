class User < ApplicationRecord
  include PrettyId
  self.id_prefix = :usr

  has_secure_password validations: false

  has_many :user_sessions, dependent: :destroy

  validates :username, uniqueness: true

  class << self
    def identify(identifier)
      return nil unless identifier.present?

      where(username: identifier).or(where(email: identifier)).first
    end
  end
end
