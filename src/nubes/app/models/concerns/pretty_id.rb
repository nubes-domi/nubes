module PrettyId
  extend ActiveSupport::Concern

  class << self
    attr_accessor :prefixes

    def find(id)
      prefix, * = id.rpartition("_")
      prefixes[prefix.to_sym].find(id)
    end
  end

  module ClassMethods
    attr_reader :id_prefix

    def id_prefix=(id_prefix)
      @id_prefix = id_prefix
      PrettyId.prefixes ||= {}
      PrettyId.prefixes[id_prefix] = self
    end
  end

  included do
    before_create do
      self.id = generate_id
    end
  end

  private

  # With an alphabet of 62 symbols (lowercase, uppercase, digits) we can reach
  # the uniqueness of UUIDs (2^124 - 4 bits are not random) with just 21 digits.
  def generate_id
    alphabet = ("0".."9").to_a + ("a".."z").to_a + ("A".."Z").to_a
    id = SecureRandom.random_bytes(21).chars.map(&:ord).map { |n| alphabet[n % alphabet.length] }.join

    "#{self.class.id_prefix}_#{id}"
  end
end
