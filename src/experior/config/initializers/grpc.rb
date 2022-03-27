require "sum_services_pb"

SUM_SESSIONS = Sum::Sessions::Stub.new("localhost:8081", :this_channel_is_insecure)
SUM_USERS = Sum::Users::Stub.new("localhost:8081", :this_channel_is_insecure)
