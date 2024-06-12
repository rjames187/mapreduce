seq_wc: build
	./mapreduce -r sequential -d ./test_text/ -p wc > seq_wc.txt

build: 
	go build

clean:
	rm mapreduce.exe
	rm seq_wc.txt