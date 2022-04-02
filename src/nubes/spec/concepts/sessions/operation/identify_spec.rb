require "rails_helper"

RSpec.describe Sessions::Operation::Identify do
  fixtures :users

  it "finds existing users" do
    result = described_class.wtf?(params: { identifier: users(:joe).username })
    expect(result).to be_success
  end

  it "fails on unknown identifiers" do
    result = described_class.wtf?(params: { identifier: "i-do-not-exist" })
    expect(result).to be_failure
  end
end
