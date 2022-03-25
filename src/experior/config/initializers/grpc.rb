require "sum_services_pb"

SUM_CLIENT = Sum::Sessions::Stub.new("localhost:8081", :this_channel_is_insecure)
