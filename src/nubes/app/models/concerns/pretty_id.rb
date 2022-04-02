module PrettyId
  extend ActiveSupport::Concern

  module ClassMethods
    attr_accessor :id_prefix
  end

  included do
    before_create do
      self.id = generate_id
    end
  end

  private

  def generate_id
    alphabet = ("0".."9").to_a + ("a".."z").to_a + ("A".."Z").to_a
    id = 16.times.map { alphabet.sample }.join

    "#{self.class.id_prefix}_#{id}"
  end
end
