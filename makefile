seq_wc: build
	./mapreduce -r sequential -d ./mock_fs/ -p wc > seq_wc.txt

build: 
	go build

clean:
	rm mapreduce.exe &
	rm seq_wc.txt &
	rm ./mock_fs/m*.txt &
	rm ./mock_fs/o*.txt