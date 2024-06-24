# start the master and create a file called 'done' when the master exits
rm -f done
((./mapreduce -r master -d ./mock_fs/ -a 127.0.0.1:12345 -nr 3); touch done) &
sleep 1
(./mapreduce -r worker -a 127.0.0.1:12345 -p c) &
(./mapreduce -r worker -a 127.0.0.1:12345 -p c) &
(./mapreduce -r worker -a 127.0.0.1:12345 -p c) &
# continuously create workers in case all crash
# stop after the 'done' file is created
( while [ ! -f done ]
do
  echo 'new worker created!'
  ./mapreduce -r worker -a 127.0.0.1:12345 -p c
  sleep 4
done) &

while [ ! -f done ]
do
  echo 'new worker created!'
  ./mapreduce -r worker -a 127.0.0.1:12345 -p c
  sleep 4
done

wait

rm -f done

# combine and sort all output files
sort ./mock_fs/o*.txt > ./mock_fs/all.txt
# compare sequential output to parallel output
if cmp ./mock_fs/all.txt ./mock_fs/seq_wc.txt
then
  echo crash output is correct
  echo '---' c: PASS
else
  echo sequential output did not match parallel output for crash test
  echo '---' c: FAIL
fi

wait
