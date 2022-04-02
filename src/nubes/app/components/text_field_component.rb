class TextFieldComponent < ViewComponent::Base
  def initialize(name:, **attrs)
    super
    @name = name
    @id = attrs.delete(:id) || name
    @type = attrs.delete(:type) || :text
    @required = attrs.delete(:required)
    @autofocus = attrs.delete(:autofocus)
    @size = attrs.delete(:size)
    @error = attrs.delete(:error)
    @attrs = attrs
  end
end
