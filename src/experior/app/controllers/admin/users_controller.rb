module Admin
  class UsersController < ApplicationController
    before_action :load_user, only: [:show, :edit, :update, :destroy]

    def index
      @users = User.list
    end

    def new
    end

    def create
    end

    def show
    end

    def edit
    end

    def update
    end

    def destroy
    end

    private

    def load_user
      @user = User.get(params[:id])
    end
  end
end
