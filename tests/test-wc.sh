# run a mapreduce job with three workers and 3 output files
./mapreduce -r master -d ./mock_fs/ -a 127.0.0.1:12345 -nr 3 &
pid=$!
sleep 1
(./mapreduce -r worker -a 127.0.0.1:12345 -p wc)
(./mapreduce -r worker -a 127.0.0.1:12345 -p wc)
(./mapreduce -r worker -a 127.0.0.1:12345 -p wc)
wait $pid

# combine and sort all output files
sort ./mock_fs/o*.txt > ./mock_fs/all.txt
# compare sequential output to parallel output
if cmp ./mock_fs/all.txt ./mock_fs/seq_wc.txt
then
  echo word count output is correct
  echo '---' wc: PASS
else
  echo sequential output did not match parallel output for word count
  echo '---' wc: FAIL
fi

wait