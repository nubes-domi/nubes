require "rails_helper"

RSpec.describe Sessions::Operation::Start do
  fixtures :users

  let(:request) do
    double(user_agent: "nubes/0.0.1", remote_addr: "127.0.0.1")
  end

  it "succeeds with valid credentials" do
    result = described_class.wtf?(params: {
      user_id: users(:joe).id,
      password: "secret"
    }, request: request)

    expect(result).to be_success
    expect(result["session"]).not_to be_nil
  end

  it "fails on unknown identifiers" do
    result = described_class.wtf?(params: { identifier: "i-do-not-exist" })
    expect(result).to be_failure
  end
end
