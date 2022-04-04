class ContactAddressComponent < ViewComponent::Base
  def initialize(**kwargs)
    super

    @presenter = AddressPresenter.new(**kwargs)
  end
end
