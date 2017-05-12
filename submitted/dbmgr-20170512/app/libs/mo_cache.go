package libs

/*
// ManagedObject 代表一个管理对象, 注意它是一个不可变对象，任何人不要试图修改它
type ManagedObject struct {
	cache *MoCache
	models.Object
	Type  *types.ClassDefinition
	Value interface{}
}

// GetName 返回对象的名称
func (mo *ManagedObject) GetName() string {
	if mo.Value == nil {
		return mo.Name
	}
	switch o := mo.Value.(type) {
	case NetworkDevice:
		return o.DisplayName
	}
	return mo.Name
}

// GetDBObject 返回对象的数据库模型
func (mo *ManagedObject) GetDBObject() interface{} {
	if mo.Value == nil {
		return mo.Object
	}
	switch o := mo.Value.(type) {
	case *NetworkDevice:
		return o.NetworkDevice
	case *NetworkLink:
		return o.NetworkLink
	}
	return mo.Value
}

// RecordVersion 记录的版本号
func (mo *ManagedObject) RecordVersion() ds.RecordVersion {
	return ds.RecordVersion{Id: mo.Id, UpdatedAt: mo.UpdatedAt}
}

// NetworkDevice 代表一个网络管理对象，任何人不要试图修改它
type NetworkDevice struct {
	mo   *ManagedObject
	Type *types.ClassDefinition
	models.NetworkDevice
}

// NetworkLink 代表一个线路，任何人不要试图修改它
type NetworkLink struct {
	mo   *ManagedObject
	Type *types.ClassDefinition
	models.NetworkLink

	from, to atomic.Value
}

// From 返回线路入端的设备
func (link *NetworkLink) From() (*NetworkDevice, error) {
	o := link.from.Load()
	if nil != o {
		if nd, ok := o.(*NetworkDevice); ok {
			return nd, nil
		}
	}

	nd, err := link.mo.cache.GetNetworkDevice(link.FromDevice)
	if err != nil {
		return nil, err
	}

	link.from.Store(nd)
	return nd, nil
}

// To 返回线路另一端的设备
func (link *NetworkLink) To() (*NetworkDevice, error) {
	o := link.to.Load()
	if nil != o {
		if nd, ok := o.(*NetworkDevice); ok {
			return nd, nil
		}
	}

	nd, err := link.mo.cache.GetNetworkDevice(link.ToDevice)
	if err != nil {
		return nil, err
	}

	link.to.Store(nd)
	return nd, nil
}

// MoCache 管理对角的缓存
type MoCache struct {
	lifecycle      *Lifecycle
	managedElement *types.ClassDefinition
	managedObject  *types.ClassDefinition
	managedLink    *types.ClassDefinition
	networkDevice  *types.ClassDefinition
	lock           sync.RWMutex
	values         map[int64]*ManagedObject
	all_devices    []*NetworkDevice
	all_links      []*NetworkLink
}

// Refresh 刷新绶存，确保内存与数据库中的数据一致。
func (cache *MoCache) Refresh() error {
	cache.lock.RLock()
	isEmpty := len(cache.values) == 0
	cache.lock.RUnlock()
	if isEmpty {
		log.Println("[mo_cache] cache is empty, skip refresh.")
		return nil
	}

	sqlStr, sqlArgs, _ := squirrel.Select("id", "updated_at").From(models.ObjectModel.TableName).ToSql()
	rows, err := cache.lifecycle.DbRunner.Query(sqlStr, sqlArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			cache.lock.Lock()
			cache.values = nil
			cache.all_devices = nil
			cache.all_links = nil
			cache.lock.Unlock()

			log.Println("[mo_cache] database is empty, clear cache.")
			return nil
		}
		return errors.New("GetSnapshots:" + err.Error())
	}
	defer rows.Close()

	moCopies := map[int64]*ManagedObject{}
	cache.lock.RLock()
	for k, v := range cache.values {
		moCopies[k] = v
	}
	cache.lock.RUnlock()

	//var created []int64
	var updated []int64
	//var deleted []int64

	for rows.Next() {
		var version ds.RecordVersion
		if err := rows.Scan(&version.Id, &version.UpdatedAt); err != nil {
			return errors.New("ReadSnapshots:" + err.Error())
		}

		if o, ok := moCopies[version.Id]; ok {
			if !o.UpdatedAt.Equal(version.UpdatedAt) {
				updated = append(updated, version.Id)
			}
			delete(moCopies, version.Id)
		} // else {
		//	created = append(created, version.Id)
		// }
	}
	cache.lock.Lock()
	for id := range moCopies {
		delete(cache.values, id)
	}
	for _, id := range updated {
		delete(cache.values, id)
	}

	cache.all_devices = nil
	cache.all_links = nil
	cache.lock.Unlock()

	log.Println("[mo_cache] update", len(updated), ", delete", len(moCopies))
	return nil
}

func findSpec(definitions *types.TableDefinitions, name string, defSpec *types.ClassDefinition) *types.ClassDefinition {
	typeSpec := definitions.FindByUnderscoreName(name)
	if typeSpec != nil {
		return typeSpec
	}

	typeSpec = definitions.Find(name)
	if typeSpec != nil {
		return typeSpec
	}

	return defSpec
}

func (cache *MoCache) toNetworkDevice(obj *models.Object) (*ManagedObject, error) {
	nd, err := models.NetworkDeviceModel.FindByID(cache.lifecycle.DbRunner, obj.Id)
	if err != nil {
		return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(obj.Id, 10) + ":" + obj.Name + ") fail, " + err.Error())
	}

	typeSpec := findSpec(cache.lifecycle.Definitions, nd.Type, cache.networkDevice)
	if typeSpec == nil {
		return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(nd.Id, 10) + ":" + nd.DisplayName + ") fail, type(" + nd.Type + ") is unknown.")
	}

	mo := &ManagedObject{
		cache:  cache,
		Object: *obj,
		Type:   typeSpec,
	}
	mo.Value = &NetworkDevice{mo: mo,
		NetworkDevice: *nd,
		Type:          typeSpec}
	return mo, nil
}

func (cache *MoCache) toNetworkDeviceFrom(nd *models.NetworkDevice) (*ManagedObject, error) {

	typeSpec := findSpec(cache.lifecycle.Definitions, nd.Type, cache.networkDevice)
	if typeSpec == nil {
		return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(nd.Id, 10) + ":" + nd.DisplayName + ") fail, type(" + nd.Type + ") is unknown.")
	}

	mo := &ManagedObject{
		cache: cache,
		Object: models.Object{
			TableName: "tpt_network_devices",
			Id:        nd.Id,
			Type:      nd.Type,
			Name:      nd.Name,
			UpdatedAt: nd.UpdatedAt,
			CreatedAt: nd.CreatedAt,
		},
		Type: typeSpec,
	}
	mo.Value = &NetworkDevice{mo: mo,
		NetworkDevice: *nd,
		Type:          typeSpec}
	return mo, nil
}

func (cache *MoCache) toNetworkLink(obj *models.Object) (*ManagedObject, error) {
	nd, err := models.NetworkLinkModel.FindByID(cache.lifecycle.DbRunner, obj.Id)
	if err != nil {
		return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(obj.Id, 10) + ":" + obj.Name + ") fail, " + err.Error())
	}

	// typeSpec := cache.lifecycle.Definitions.FindByUnderscoreName(nd.Type)
	// if typeSpec == nil {
	// 	typeSpec = cache.lifecycle.Definitions.Find(nd.Type)
	// 	if typeSpec == nil {
	// 		return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(nd.Id, 10) + ":" + nd.DisplayName + ") fail, type is unknown.")
	// 	}
	// }
	typeSpec := cache.managedLink

	mo := &ManagedObject{
		cache:  cache,
		Object: *obj,
		Type:   typeSpec,
	}
	mo.Value = &NetworkLink{mo: mo,
		NetworkLink: *nd,
		Type:        typeSpec}
	return mo, nil
}

func (cache *MoCache) toManagedObject(obj *models.Object) (*ManagedObject, error) {
	switch obj.TableName {
	case "tpt_network_devices":
		return cache.toNetworkDevice(obj)
	case "tpt_network_links":
		return cache.toNetworkLink(obj)
	default:
		typeSpec := cache.lifecycle.Definitions.FindByTableName(obj.TableName)
		if typeSpec == nil {
			return nil, errors.New("toManagedObject: load mo(" + strconv.FormatInt(obj.Id, 10) + ":" + obj.Name + ") fail, type is unknown.")
		}
		return &ManagedObject{
			cache:  cache,
			Object: *obj,
			Type:   typeSpec,
		}, nil
	}
}

// Get 获取一个指定 ID 的管理对象，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) Get(moID int64) (*ManagedObject, error) {
	cache.lock.RLock()
	if cache.values != nil {
		if old, ok := cache.values[moID]; ok {
			cache.lock.RUnlock()
			return old, nil
		}
	}
	cache.lock.RUnlock()

	obj, err := models.ObjectModel.FindByID(cache.lifecycle.DbRunner, moID)
	if err != nil {
		return nil, err
	}
	if nil == obj {
		return nil, NotFound(moID)
	}

	mo, err := cache.toManagedObject(obj)
	if err != nil {
		return nil, err
	}
	cache.lock.Lock()
	if cache.values == nil {
		cache.values = map[int64]*ManagedObject{moID: mo}
	} else {
		cache.values[moID] = mo
	}
	cache.lock.Unlock()

	return mo, nil
}

func (cache *MoCache) get(obj *models.Object) (*ManagedObject, error) {
	cache.lock.RLock()
	if cache.values != nil {
		if old, ok := cache.values[obj.Id]; ok {
			cache.lock.RUnlock()
			return old, nil
		}
	}
	cache.lock.RUnlock()

	mo, err := cache.toManagedObject(obj)
	if err != nil {
		return nil, err
	}
	cache.lock.Lock()
	if cache.values == nil {
		cache.values = map[int64]*ManagedObject{obj.Id: mo}
	} else {
		cache.values[obj.Id] = mo
	}
	cache.lock.Unlock()

	return mo, nil
}

// GetNetworkDevice 获取一个指定 ID 的网络管理对象，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) GetNetworkDevice(moID int64) (*NetworkDevice, error) {
	mo, err := cache.Get(moID)
	if err != nil {
		return nil, err
	}
	if nil == mo {
		return nil, NotFound(moID)
	}
	nd, ok := mo.Value.(*NetworkDevice)
	if !ok {
		return nil, errors.New("GetNetworkDevice: load mo(" + strconv.FormatInt(mo.Id, 10) + ":" + mo.Name + ") fail, type(" + mo.Type.UName() + ") isn't network device.")
	}
	return nd, nil
}

func (cache *MoCache) searchNetworkDeviceByAddress(domain, address string) ([]*models.NetworkDevice, error) {
	filters := []models.Sqlizer{models.NetworkDeviceModel.C.ADDRESS.EQU(address)}
	if "" != domain {
		subFilter := models.DomainModel.Where(models.DomainModel.C.NAME.EQU(domain),
			models.NetworkDeviceModel.C.DOMAINID.
				TableAlias(models.NetworkDeviceModel.TableName).
				EQU(models.DomainModel.C.ID.Name)).
			Select(models.DomainModel.C.ID.Name).
			From(models.DomainModel.TableName)
		filters = append(filters, models.EXISTS(subFilter))
	}

	builder := models.NetworkDeviceModel.Where(filters...).Select()
	devices, err := models.NetworkDeviceModel.QueryWith(cache.lifecycle.DbRunner, builder)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		addressFilter := models.NetworkAddressModel.Where(models.NetworkAddressModel.C.ADDRESS.EQU(address),
			models.NetworkDeviceModel.C.ID.TableAlias(models.NetworkAddressModel.TableName).
				EQU(models.NetworkAddressModel.C.MANAGEDOBJECTID.Name)).
			Select(models.NetworkAddressModel.C.MANAGEDOBJECTID.Name).
			From(models.DomainModel.TableName)
		filters := []models.Sqlizer{models.EXISTS(addressFilter)}
		if "" != domain {
			subFilter := models.DomainModel.Where(models.DomainModel.C.NAME.EQU(domain),
				models.NetworkDeviceModel.C.DOMAINID.
					TableAlias(models.NetworkDeviceModel.TableName).
					EQU(models.DomainModel.C.ID.Name)).
				Select(models.DomainModel.C.ID.Name).
				From(models.DomainModel.TableName)
			filters = append(filters, models.EXISTS(subFilter))
		}

		builder = models.NetworkDeviceModel.Where(filters...).Select()
		devices, err = models.NetworkDeviceModel.QueryWith(cache.lifecycle.DbRunner, builder)
		if err != nil {
			return nil, err
		}
	}
	return devices, nil
}

// GetNetworkDeviceByAddress 获取一个指定 ID 的网络管理对象，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) GetNetworkDeviceByAddress(domain, address string) (*NetworkDevice, error) {
	devices, err := cache.searchNetworkDeviceByAddress(domain, address)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, nil
	}
	if len(devices) != 1 {
		return nil, errors.New("GetNetworkDeviceByAddress: muti choice is find.")
	}
	return cache.loadOrCreateNetworkDevice(devices[0])
}

// SearchNetworkDeviceByName 获取一个指定 ID 的网络管理对象，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) SearchNetworkDeviceByName(name string) ([]*NetworkDevice, error) {
	filter := models.NetworkDeviceModel.C.NAME.LIKE(name)
	builder := models.NetworkDeviceModel.Where(filter).Select()
	devices, err := models.NetworkDeviceModel.QueryWith(cache.lifecycle.DbRunner, builder)
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		filter = models.NetworkDeviceModel.C.ZHNAME.LIKE(name)
		builder = models.NetworkDeviceModel.Where(filter).Select()
		devices, err := models.NetworkDeviceModel.QueryWith(cache.lifecycle.DbRunner, builder)
		if err != nil {
			return nil, err
		}
		if len(devices) == 0 {
			return nil, nil
		}
	}
	return cache.convertNetworkDevices(devices)
}

func (cache *MoCache) convertNetworkDevices(devices []*models.NetworkDevice) ([]*NetworkDevice, error) {
	results := make([]*NetworkDevice, 0, len(devices))
	for _, dev := range devices {
		nd, err := cache.loadOrCreateNetworkDevice(dev)
		if err != nil {
			return nil, err
		}
		results = append(results, nd)
	}
	return results, nil
}

// SearchNetworkDevices 按指定的字符在设备的名称或地址字段中查找，注意名称是模糊查找
func (cache *MoCache) SearchNetworkDevices(name string) ([]*NetworkDevice, error) {
	if ip := net.ParseIP(name); nil != ip {
		devices, err := cache.searchNetworkDeviceByAddress("", name)
		if err != nil {
			return nil, err
		}
		if len(devices) == 0 {
			return nil, nil
		}
		return cache.convertNetworkDevices(devices)
	}
	return cache.SearchNetworkDeviceByName(name)
}

func (cache *MoCache) loadOrCreateNetworkDevice(nd *models.NetworkDevice) (*NetworkDevice, error) {
	cache.lock.RLock()
	if cache.values != nil {
		if old, ok := cache.values[nd.Id]; ok {
			cache.lock.RUnlock()

			nd, ok := old.Value.(*NetworkDevice)
			if !ok {
				return nil, errors.New("loadOrCreateNetworkDevice: load mo(" + strconv.FormatInt(nd.Id, 10) +
					":" + nd.Name + ") fail, type(" + nd.Type.UName() + ") isn't network device.")
			}
			return nd, nil
		}
	}
	cache.lock.RUnlock()

	mo, err := cache.toNetworkDeviceFrom(nd)
	if err != nil {
		return nil, err
	}
	cache.lock.Lock()
	if cache.values == nil {
		cache.values = map[int64]*ManagedObject{mo.Id: mo}
	} else {
		cache.values[mo.Id] = mo
	}
	cache.lock.Unlock()

	return mo.Value.(*NetworkDevice), nil
}

// ListNetworkDevices 列出所有的网络管理对象，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) ListNetworkDevices() ([]*NetworkDevice, error) {
	cache.lock.RLock()
	all := cache.all_devices
	cache.lock.RUnlock()
	if nil != all {
		return all, nil
	}

	moList, err := models.ObjectModel.QueryWith(cache.lifecycle.DbRunner,
		models.ObjectModel.Where(models.ObjectModel.C.TABLENAME.EQU("tpt_network_devices")).Select())
	if nil != err {
		return nil, err
	}

	var devices = make([]*NetworkDevice, 0, len(moList))
	for _, mo := range moList {
		mo, err := cache.get(mo)
		if err != nil {
			return nil, err
		}

		nd, ok := mo.Value.(*NetworkDevice)
		if !ok {
			return nil, errors.New("ListNetworkDevices: load mo(" + strconv.FormatInt(mo.Id, 10) + ":" + mo.Name + ") fail, type(" + mo.Type.UName() + ") isn't network device.")
		}
		devices = append(devices, nd)
	}

	cache.lock.Lock()
	cache.all_devices = devices
	cache.lock.Unlock()
	return devices, nil
}

// GetNetworkLink 获取一个指定 ID 的网络线路，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) GetNetworkLink(moID int64) (*NetworkLink, error) {
	mo, err := cache.Get(moID)
	if err != nil {
		return nil, err
	}
	if mo == nil {
		return nil, NotFound(moID)
	}

	nd, ok := mo.Value.(*NetworkLink)
	if !ok {
		return nil, errors.New("GetNetworkLink: load mo(" + strconv.FormatInt(mo.Id, 10) + ":" + mo.Name + ") fail, type(" + mo.Type.UName() + ") isn't network link.")
	}
	return nd, nil
}

// ListNetworkLinks 列出所有的网络线路，如果管理对象没有被加功到内存那么立即加载, 注意它是一个不可变对象，任何人不要试图修改它
func (cache *MoCache) ListNetworkLinks() ([]*NetworkLink, error) {
	cache.lock.RLock()
	all := cache.all_links
	cache.lock.RUnlock()
	if nil != all {
		return all, nil
	}

	moList, err := models.ObjectModel.QueryWith(cache.lifecycle.DbRunner,
		models.ObjectModel.Where(models.ObjectModel.C.TABLENAME.EQU("tpt_network_links")).Select())
	if nil != err {
		return nil, err
	}

	var links = make([]*NetworkLink, 0, len(moList))
	for _, mo := range moList {
		mo, err := cache.get(mo)
		if err != nil {
			return nil, err
		}

		nd, ok := mo.Value.(*NetworkLink)
		if !ok {
			return nil, errors.New("ListNetworkLinks: load mo(" + strconv.FormatInt(mo.Id, 10) + ":" + mo.Name + ") fail, type(" + mo.Type.UName() + ") isn't network link.")
		}
		links = append(links, nd)
	}

	cache.lock.Lock()
	cache.all_links = links
	cache.lock.Unlock()
	return links, nil
}
*/
