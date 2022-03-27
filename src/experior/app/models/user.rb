class User
  include ActiveModel::API

  attr_accessor :id, :name, :username, :email, :admin, :picture

  def list_name
    name.presence || username
  end

  def to_param
    id
  end

  class << self
    def from_grpc(model)
      new(
        id: model.id,
        name: model.name,
        username: model.username,
        email: model.email,
        admin: model.admin,
        picture: model.picture
      )
    end

    def list
      result = SUM_USERS.list(Sum::ListUsersRequest.new, { metadata: { authorization: Thread.current[:authorization] }})

      result.users.map { |u| from_grpc(u) }
    end

    def get(id)
      from_grpc SUM_USERS.get(Sum::GetUserRequest.new(user_id: id), { metadata: { authorization: Thread.current[:authorization] }})
    rescue GRPC::Unauthenticated
      false
    end
  end
end
