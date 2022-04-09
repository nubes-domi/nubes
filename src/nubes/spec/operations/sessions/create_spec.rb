require "rails_helper"

RSpec.describe Sessions::Create do
  fixtures :users

  let(:request) do
    double(user_agent: "nubes/0.0.1", remote_addr: "127.0.0.1")
  end

  it "succeeds with valid credentials" do
    result = subject.call(
      user_id: users(:joe).id,
      password: "secret",
      request:
    )

    expect(result).to be_success
    expect(result.value!).to be_a UserSession
  end

  it "fails on wrong password" do
    result = subject.call(
      user_id: users(:joe).id,
      password: "wrong-password",
      request:
    )
    expect(result).to be_failure
  end
end
