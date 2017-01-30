<?php

class ddb {
	
	private $host;
	private $port;
	private $socket;
	
	public function __construct(string $host, int $port){
		$this->host = $host;
		$this->port = $port;
	}
	
	public function insert($target, $data){
		$writeCommand = [
			'command' => 'write',
			'target' => $target,
			'data' => $data
		];
		
		return $this->sendCommandToSocket($writeCommand);
	}
	
	public function select($target, $whereData){
		$readCommand = [
			'command' => 'read',
			'target' => $target,
			'data' => $whereData
		];
				
		return $this->sendCommandToSocket($readCommand);
	}
	
	public function close(){
		$exitCommand = ['command' => 'exit'];
		$this->sendCommandToSocket($exitCommand);
		$socket = $this->getSocket(true);
		if($socket !== null){
			fclose($this->socket);
			$this->socket = null;
		}
		return true;
	}
	
	private function sendCommandToSocket($command){
		$socket = $this->getSocket();
		fwrite($socket, json_encode($command) . "\r\n");
		$result = json_decode(trim(fgets($socket)), true);
		
		return $result;
	}
	
	private function getSocket($forClose = false){
		if($this->socket === null){
			if($forClose === false){
				$this->socket = fsockopen($this->host, $this->port);
			}
		}
		
		return $this->socket;
	}
}

$ddb = new ddb('localhost', 9977);

$insertResult = $ddb->insert('foobar', ['foo' => 'bar', 'bar' => 'foo']);

if(array_key_exists('_id', $insertResult)){
	$selectData = $ddb->select('foobar', ['_id' => $insertResult['_id']]);
	var_dump($selectData);
}

$ddb->close();