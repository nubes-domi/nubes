class TextFieldComponent < ViewComponent::Base
  def initialize(name:, id: nil, type: :text, required: false, autofocus: false, error: nil, size: nil, **attrs)
    super
    @name = name
    @id = id || name
    @type = type
    @required = required
    @autofocus = autofocus
    @size = size
    @error = error
    @attrs = attrs
  end
end
