package models

import (
	walletPB "github.com/intellisoftalpin/proto/proto-gen/wallet"
)

type NetworkInfoResponse struct {
	NetworkInfo  NetworkInfo  `json:"network_info"`
	NetworkTip   NetworkTip   `json:"network_tip"`
	NextEpoch    NextEpoch    `json:"next_epoch"`
	NodeEra      string       `json:"node_era"`
	NodeTip      NodeTip      `json:"node_tip"`
	SyncProgress SyncProgress `json:"sync_progress"`
	WalletMode   string       `json:"wallet_mode"`
}

type NetworkInfo struct {
	NetworkID     string `json:"network_id"`
	ProtocolMagic uint64 `json:"protocol_magic"`
}

type NetworkTip struct {
	AbsoluteSlotNumber uint64 `json:"absolute_slot_number"`
	EpochNumber        uint64 `json:"epoch_number"`
	SlotNumber         uint64 `json:"slot_number"`
	Time               string `json:"time"`
}

type NextEpoch struct {
	EpochNumber    uint64 `json:"epoch_number"`
	EpochStartTime string `json:"epoch_start_time"`
}

type NodeTip struct {
	AbsoluteSlotNumber uint64   `json:"absolute_slot_number"`
	EpochNumber        uint64   `json:"epoch_number"`
	Height             Quantity `json:"height"`
	SlotNumber         uint64   `json:"slot_number"`
	Time               string   `json:"time"`
}

type SyncProgress struct {
	Status   string   `json:"status"`
	Progress Quantity `json:"progress"`
}

type Quantity struct {
	Quantity uint64 `json:"quantity"`
	Unit     string `json:"unit"`
}

func ToNetworkInfoResponse(nInfo *walletPB.GetWalletNetworkInfoResponse) NetworkInfoResponse {
	return NetworkInfoResponse{
		NetworkInfo: NetworkInfo{
			NetworkID:     nInfo.NetworkInfo.NetworkId,
			ProtocolMagic: nInfo.NetworkInfo.ProtocolMagic,
		},
		NetworkTip: NetworkTip{
			AbsoluteSlotNumber: nInfo.NetworkTip.AbsoluteSlotNumber,
			EpochNumber:        nInfo.NetworkTip.EpochNumber,
			SlotNumber:         nInfo.NetworkTip.SlotNumber,
			Time:               nInfo.NetworkTip.Time,
		},
		NextEpoch: NextEpoch{
			EpochNumber:    nInfo.NextEpoch.EpochNumber,
			EpochStartTime: nInfo.NextEpoch.EpochStartTime,
		},
		NodeEra: nInfo.NodeEra,
		NodeTip: NodeTip{
			AbsoluteSlotNumber: nInfo.NodeTip.AbsoluteSlotNumber,
			EpochNumber:        nInfo.NodeTip.EpochNumber,
			Height: Quantity{
				Quantity: nInfo.NodeTip.Height.Quantity,
				Unit:     nInfo.NodeTip.Height.Unit,
			},
			SlotNumber: nInfo.NodeTip.SlotNumber,
			Time:       nInfo.NodeTip.Time,
		},
		SyncProgress: SyncProgress{
			Status: nInfo.SyncProgress.Status,
			Progress: Quantity{
				Quantity: nInfo.SyncProgress.Progress.Quantity,
				Unit:     nInfo.SyncProgress.Progress.Unit,
			},
		},
		WalletMode: nInfo.WalletMode,
	}
}
