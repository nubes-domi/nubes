require "rails_helper"

RSpec.describe UserSession, type: :model do
  fixtures :users

  it "uses pretty ids" do
    session = UserSession.create(user: users(:joe))
    expect(session.id).to start_with "usr_sess_"
  end
end
