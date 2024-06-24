test_wc: clean seq_wc
	./tests/test-wc.sh

test_fault: clean seq_wc
	./tests/test-fault.sh

seq_wc: build
	./mapreduce -r sequential -d ./mock_fs/ -p wc > mock_fs/seq_wc_int.txt
	cd mock_fs
	sort mock_fs/seq_wc_int.txt > mock_fs/seq_wc.txt

build: 
	go build

clean:
	rm -f mapreduce.exe &
	rm -f seq_wc.txt &
	rm -f ./mock_fs/m*.txt &
	rm -f ./mock_fs/o*.txt &
	rm -f ./mock_fs/seq*.txt &
	rm -f ./mock_fs/seq*int.txt &
	rm -f ./mock_fs/all.txt
	rm -f done