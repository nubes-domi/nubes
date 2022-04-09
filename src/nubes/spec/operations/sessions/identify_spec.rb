require "rails_helper"

RSpec.describe Sessions::Identify do
  fixtures :users

  it "finds existing users" do
    result = subject.call(identifier: users(:joe).username)
    expect(result).to be_success
  end

  it "fails on unknown identifiers" do
    result = subject.call(identifier: "i-do-not-exist")
    expect(result).to be_failure
  end
end
