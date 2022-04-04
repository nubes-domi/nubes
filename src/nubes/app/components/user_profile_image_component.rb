# frozen_string_literal: true

class UserProfileImageComponent < ViewComponent::Base
  def initialize(user:, size: "md", klass: "", **kwargs)
    super
    @user = user
    @size = size
    @class = klass
    @attrs = kwargs
  end

  def picture?
    false #@user.picture.present?
  end

  def placeholder_color
    colors = %w[659a10 109a39 109a7b 107f9a 104f9a 55109a 8f109a 979a10 9a1110 9a6410]

    n = Zlib.crc32(@user.id) % colors.length

    "##{colors[n]}"
  end

  def placeholder_letters
    n = @user.list_name
    case n
    when / /
      n.split.first(2).map { |part| part[0] }.join
    when /\./
      n.split(".").first(2).map { |part| part[0] }.join
    else
      n.first(2)
    end.upcase
  end
end
