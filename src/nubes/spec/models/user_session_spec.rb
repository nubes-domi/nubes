require "rails_helper"

RSpec.describe UserSession, type: :model do
  fixtures :users, :user_sessions

  it "uses pretty ids" do
    session = UserSession.create(user: users(:joe))
    expect(session.id).to start_with "usr_sess_"
  end

  describe "#token" do
    it "generates a JWT that can be used to retrieve itself" do
      session = user_sessions(:joe1)
      token = session.token

      expect(token).to be_present
      expect(UserSession.for_token(token)).to eq session
    end
  end

  describe ".for_token" do
    it "rejects invalid tokens" do
      expect(UserSession.for_token("bad")).to be_nil
    end
  end
end
