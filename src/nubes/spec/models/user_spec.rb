require "rails_helper"

RSpec.describe User, type: :model do
  fixtures :users

  it "uses pretty ids" do
    user = User.create(username: "john")
    expect(user.id).to start_with "usr_"
  end
end
