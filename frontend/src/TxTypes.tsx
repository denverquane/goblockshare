export interface SetAliasTrans {
  Alias: string;
}

export interface LayerRepTrans {
  TxID: string;
  LayerHash: string;
  WasLayerReceived: boolean;
  WasLayerValid: boolean;
}

export interface SharedLayerTrans {
  SharedLayerHash: string;
  Recipient:       string;
}

export interface RepMessage {
  WasValid:     boolean;
  HighQuality:  boolean;
  AccurateName: boolean;
}

export interface TorrentRepTrans {
  TxID:        string;
  TorrentHash: string;
  RepMessage:  RepMessage;
}

export interface TorrentFile {
  Name:          string;
  LayerByteSize: number;
  TotalByteSize: number;
  TotalHash:     string;

  LayerHashKeys: string[];
}

export interface PublishTorrentTrans {
  Torrent: TorrentFile;
}