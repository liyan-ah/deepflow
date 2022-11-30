/*
 * Copyright (c) 2022 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package updater

import (
	cloudmodel "github.com/deepflowys/deepflow/server/controller/cloud/model"
	"github.com/deepflowys/deepflow/server/controller/db/mysql"
	"github.com/deepflowys/deepflow/server/controller/recorder/cache"
	"github.com/deepflowys/deepflow/server/controller/recorder/common"
	"github.com/deepflowys/deepflow/server/controller/recorder/db"
)

type LANIP struct {
	UpdaterBase[cloudmodel.IP, mysql.LANIP, *cache.LANIP]
}

func NewLANIP(wholeCache *cache.Cache) *LANIP {
	updater := &LANIP{
		UpdaterBase[cloudmodel.IP, mysql.LANIP, *cache.LANIP]{
			cache:        wholeCache,
			dbOperator:   db.NewLANIP(),
			diffBaseData: wholeCache.LANIPs,
		},
	}
	updater.dataGenerator = updater
	updater.cacheHandler = updater
	return updater
}

func (i *LANIP) SetCloudData(cloudData []cloudmodel.IP) {
	i.cloudData = cloudData
}

func (i *LANIP) getDiffBaseByCloudItem(cloudItem *cloudmodel.IP) (diffBase *cache.LANIP, exists bool) {
	diffBase, exists = i.diffBaseData[cloudItem.Lcuuid]
	return
}

func (i *LANIP) generateDBItemToAdd(cloudItem *cloudmodel.IP) (*mysql.LANIP, bool) {
	vinterfaceID, exists := i.cache.GetVInterfaceIDByLcuuid(cloudItem.VInterfaceLcuuid)
	if !exists {
		log.Error(resourceAForResourceBNotFound(
			common.RESOURCE_TYPE_VINTERFACE_EN, cloudItem.VInterfaceLcuuid,
			common.RESOURCE_TYPE_LAN_IP_EN, cloudItem.Lcuuid,
		))
		return nil, false
	}
	networkID, exists := i.cache.GetNetworkIDByVInterfaceLcuuid(cloudItem.VInterfaceLcuuid)
	if !exists {
		log.Error(resourceAForResourceBNotFound(
			common.RESOURCE_TYPE_VINTERFACE_EN, cloudItem.VInterfaceLcuuid,
			common.RESOURCE_TYPE_LAN_IP_EN, cloudItem.Lcuuid,
		))
		return nil, false
	}
	ip := common.FormatIP(cloudItem.IP)
	if ip == "" {
		log.Error(ipIsInvalid(
			common.RESOURCE_TYPE_LAN_IP_EN, cloudItem.Lcuuid, cloudItem.IP,
		))
		return nil, false
	}
	dbItem := &mysql.LANIP{
		IP:           ip,
		Domain:       i.cache.DomainLcuuid,
		SubDomain:    cloudItem.SubDomainLcuuid,
		NetworkID:    networkID,
		VInterfaceID: vinterfaceID,
	}
	dbItem.Lcuuid = cloudItem.Lcuuid
	return dbItem, true
}

// 保留接口
func (i *LANIP) generateUpdateInfo(diffBase *cache.LANIP, cloudItem *cloudmodel.IP) (map[string]interface{}, bool) {
	return nil, false
}

func (i *LANIP) addCache(dbItems []*mysql.LANIP) {
	i.cache.AddLANIPs(dbItems)
}

// 保留接口
func (i *LANIP) updateCache(cloudItem *cloudmodel.IP, diffBase *cache.LANIP) {
}

func (i *LANIP) deleteCache(lcuuids []string) {
	i.cache.DeleteLANIPs(lcuuids)
}
