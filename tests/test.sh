./mapreduce -r master -d ./mock_fs/ -a 127.0.0.1:12345 &
(sleep 1
./mapreduce -r worker -a 127.0.0.1:12345
./mapreduce -r worker -a 127.0.0.1:12345)