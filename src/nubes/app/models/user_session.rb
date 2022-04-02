class UserSession < ApplicationRecord
  include PrettyId
  self.id_prefix = :usr_sess

  belongs_to :user

  before_create do
    self.expires_at ||= 10.years.from_now
  end
end
