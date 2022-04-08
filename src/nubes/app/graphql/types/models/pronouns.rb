module Types
  module Models
    class Pronouns < Types::Base::Enum
      value "MASCULINE", "Refer to this person using he / him / his"
      value "FEMININE", "Refer to this person using she / her / hers"
      value "NEUTRAL", "Refer to this person using they / them / theirs"
    end
  end
end
