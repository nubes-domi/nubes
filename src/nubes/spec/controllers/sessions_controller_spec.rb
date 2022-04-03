require "rails_helper"

RSpec.describe SessionsController, type: :controller do
  fixtures :users, :user_sessions

  describe "#new" do
    it "succeeds" do
      get :new
      expect(response).to be_successful
    end
  end

  let(:user) { users(:joe) }

  describe "#create" do
    context "with valid identifier" do
      it "proceeds to credentials" do
        post :create, params: { identifier: user.username }
        expect(response).to redirect_to("/signin/password?user_id=#{user.id}")
      end
    end

    context "with invalid identifier" do
      it "redirects back" do
        post :create, params: { identifier: "invalid" }
        expect(response).to redirect_to("/signin")
      end
    end
  end

  describe "#show" do
    context "with valid user id" do
      it "succeeds" do
        get :show, params: { method: :password, user_id: user.id }
        expect(response).to be_successful
      end
    end

    context "with invalid user id" do
      it "fails" do
        expect {
          get :show, params: { method: :password, user_id: "invalid" }
        }.to raise_error(ActiveRecord::RecordNotFound)
      end
    end

    context "when a session for the user already exists in the browser" do
      let(:session) { user_sessions(:joe1) }
      let(:other_session) { user_sessions(:jane1) }

      before do
        cookies[:sessions] = "#{session.token}|#{other_session.token}"
      end

      context "when it is the current session" do
        before do
          cookies[:current_session] = session.token
        end

        it "completes authentication" do
          get :show, params: { method: :password, user_id: user.id }
          expect(response).to redirect_to("/")
        end
      end

      context "when it is not the current session" do
        before do
          cookies[:current_session] = other_session.token
        end

        it "completes authentication" do
          get :show, params: { method: :password, user_id: user.id }
          expect(response).to redirect_to("/")
        end

        it "switches the session" do
          get :show, params: { method: :password, user_id: user.id }
          expect(cookies[:current_session]).to eq session.token
        end
      end
    end
  end

  describe "#update" do
    context "with invalid password" do
      it "asks for retry" do
        post :update, params: { method: :password, user_id: user.id, password: "meh" }
        expect(response).to redirect_to("/signin/password?user_id=#{user.id}")
      end
    end

    context "with valid user id and password" do
      it "completes authentication" do
        post :update, params: { method: :password, user_id: user.id, password: "secret" }
        expect(response).to redirect_to("/")

        expect(cookies[:current_session]).not_to be_nil
        expect(cookies[:sessions]).not_to be_nil
      end

      context "when another session exists for another user" do
        let(:other_session) { user_sessions(:jane1) }

        before do
          cookies[:sessions] = other_session.token
          cookies[:current_session] = other_session.token
        end

        it "completes authentication" do
          post :update, params: { method: :password, user_id: user.id, password: "secret" }
          expect(response).to redirect_to("/")
        end

        it "switches the session" do
          post :update, params: { method: :password, user_id: user.id, password: "secret" }
          expect(cookies[:sessions].split("|").length).to eq 2
          expect(cookies[:current_session]).not_to eq other_session.token
        end

        it "preserves the existing session" do
          post :update, params: { method: :password, user_id: user.id, password: "secret" }
          expect(cookies[:sessions].split("|")).to include other_session.token
        end
      end
    end
  end
end
