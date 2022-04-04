class AddressPresenter
  attr_reader :country, :postal_code, :region, :locality, :street

  def initialize(country: nil, postal_code: nil, region: nil, locality: nil, street: nil)
    @country = country
    @postal_code = postal_code
    @region = region
    @locality = locality
    @street = street
  end

  def lines
    lines = format.map do |line|
      Array(line).map { |component| send(component) }.compact.join(" ")
    end

    lines << @country.upcase if @country.present?

    lines.select(&:present?)
  end

  def to_s
    lines.join("\n")
  end

  private

  def format
    country_key = @country&.downcase&.to_sym
    ADDRESS_FORMATS[country_key] || ADDRESS_FORMATS[:other]
  end

  ADDRESS_FORMATS = {
    gb: [
      :street,
      :locality,
      :region,
      :postal_code
    ],
    it: [
      :street,
      [:postal_code, :locality, :region]
    ],
    us: [
      :street,
      [:locality, :region, :postal_code]
    ],
    other: [
      :street,
      [:postal_code, :locality]
    ]
  }.freeze
end
